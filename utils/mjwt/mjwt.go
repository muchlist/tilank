package mjwt

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"tilank/utils/logger"
	"tilank/utils/rest_err"
	"time"
)

var (
	JwtObj JWTAssumer
	secret []byte
)

func NewJwt() JWTAssumer {
	return JwtObj
}

func init() {
	secret = []byte(os.Getenv(secretKey))
	if string(secret) == "" {
		secret = []byte("rahasia")
	}

	JwtObj = &jwtUtils{}
}

type JWTAssumer interface {
	GenerateToken(claims CustomClaim) (string, resterr.APIError)
	ValidateToken(tokenString string) (*jwt.Token, resterr.APIError)
	ReadToken(token *jwt.Token) (*CustomClaim, resterr.APIError)
}

type jwtUtils struct {
}

const (
	CLAIMS    = "claims"
	secretKey = "SECRET_KEY"

	identityKey  = "identity"
	nameKey      = "name"
	rolesKey     = "roles"
	branchKey    = "branch"
	tokenTypeKey = "type"
	expKey       = "exp"
	freshKey     = "fresh"
)

// GenerateToken membuat token jwt untuk login header, untuk menguji nilai payloadnya
// dapat menggunakan situs jwt.io
func (j *jwtUtils) GenerateToken(claims CustomClaim) (string, resterr.APIError) {
	expired := time.Now().Add(time.Minute * claims.ExtraMinute).Unix()

	jwtClaim := jwt.MapClaims{}
	jwtClaim[identityKey] = claims.Identity
	jwtClaim[nameKey] = claims.Name
	jwtClaim[rolesKey] = claims.Roles
	jwtClaim[branchKey] = claims.Branch
	jwtClaim[expKey] = expired
	jwtClaim[tokenTypeKey] = claims.Type
	jwtClaim[freshKey] = claims.Fresh

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		logger.Error("gagal menandatangani token", err)
		return "", resterr.NewInternalServerError("gagal menandatangani token", err)
	}

	return signedToken, nil
}

// ReadToken membaca inputan token dan menghasilkan pointer struct CustomClaim
// struct CustomClaim digunakan untuk nilai passing antar middleware
func (j *jwtUtils) ReadToken(token *jwt.Token) (*CustomClaim, resterr.APIError) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Error("gagal mapping token atau token tidak valid", nil)
		return nil, resterr.NewInternalServerError("gagal mapping token", nil)
	}

	customClaim := CustomClaim{
		Identity: claims[identityKey].(string),
		Name:     claims[nameKey].(string),
		Exp:      int64(claims[expKey].(float64)),
		Roles:    iToSliceString(claims[rolesKey]),
		Branch:   claims[branchKey].(string),
		Type:     int(claims[tokenTypeKey].(float64)),
		Fresh:    claims[freshKey].(bool),
	}

	return &customClaim, nil
}

func iToSliceString(assumedSliceInterface interface{}) []string {
	sliceInterface := assumedSliceInterface.([]interface{})
	sliceString := make([]string, len(sliceInterface))
	for i, v := range sliceInterface {
		sliceString[i] = v.(string)
	}

	return sliceString
}

// ValidateToken memvalidasi apakah token string masukan valid, termasuk memvalidasi apabila field exp nya kadaluarsa
func (j *jwtUtils) ValidateToken(tokenString string) (*jwt.Token, resterr.APIError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, resterr.NewAPIError("Token signing method salah", http.StatusUnprocessableEntity, "jwt_error", nil)
		}
		return secret, nil
	})

	// Jika expired akan muncul disini asalkan ada claims exp
	if err != nil {
		return nil, resterr.NewAPIError("Token tidak valid", http.StatusUnprocessableEntity, "jwt_error", []interface{}{err.Error()})
	}

	return token, nil
}
