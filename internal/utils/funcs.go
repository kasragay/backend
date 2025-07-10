package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
)

var salt = os.Getenv("PASSWORD_HASH_SALT")
var cost_ = os.Getenv("PASSWORD_HASH_COST")
var cost int

func init() {
	var err error
	cost, err = strconv.Atoi(cost_)
	if err != nil {
		cost = bcrypt.DefaultCost
	}
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password, givenPassword string) (ok bool) {
	if givenPassword == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(givenPassword+salt))
	return err == nil
}

func GetenvAsMinuteDuration(key string, defaultValue time.Duration, required bool) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" && required {
		return defaultValue, fmt.Errorf("environment variable %s is required", key)
	}
	if val == "" {
		return defaultValue, nil
	}
	minutes, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue, fmt.Errorf("environment variable %s must be an integer", key)
	}
	return time.Duration(minutes) * time.Minute, nil
}

func ImageReader(imgData []byte, convertTo string, size [2]int, validFormats []string) (img_ *image.Image, err error) {
	defer func() {
		const caller = packageCaller + ".ImageReader"
		if err != nil {
			var uErr *Error
			if errors.As(err, &uErr) {
				err = uErr.WithCaller(caller)
			}
		}
	}()
	const maxInputBytes = 10 * 1024 * 1024
	const maxWidth, maxHeight = 5000, 5000
	if (size[0] == 0 && size[1] != 0) || (size[0] != 0 && size[1] == 0) || size[0] < 0 || size[1] < 0 {
		return nil, fmt.Errorf("invalid size: %v", size)
	}
	if len(validFormats) == 0 {
		return nil, fmt.Errorf("invalid valid formats: %v", validFormats)
	}

	if len(imgData) > maxInputBytes {
		return nil, BadAvatarResponse.Clone().
			WithReason("size", len(imgData)).
			WithReason("max_size", maxInputBytes)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return nil, BadAvatarResponse.Clone().
			WithReason("error", "error decoding image")
	}
	if cfg.Width > maxWidth || cfg.Height > maxHeight {
		return nil, BadAvatarResponse.Clone().
			WithReason("width", cfg.Width).
			WithReason("height", cfg.Height).
			WithReason("max_width", maxWidth).
			WithReason("max_height", maxHeight)
	}
	if !slices.Contains(validFormats, format) {
		return nil, BadAvatarResponse.Clone().
			WithReason("format", format).
			WithReason("valid_formats", validFormats)
	}
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, BadAvatarResponse.Clone().
			WithReason("error", "error decoding image")
	}
	if size[0] != 0 && size[1] != 0 && (cfg.Width != size[0] || cfg.Height != size[1]) {
		img = resize.Resize(uint(size[0]), uint(size[1]), img, resize.Lanczos3)
	}
	var buf bytes.Buffer
	switch convertTo {
	case "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, img)
	default:
		return &img, nil
	}
	if err != nil {
		return nil, BadAvatarResponse.Clone().
			WithReason("error", "error encoding image")
	}

	img, _, err = image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, BadAvatarResponse.Clone().
			WithReason("error", "error decoding image")
	}
	return &img, nil
}

func StructToMap(obj any) (map[string]any, error) {
	result := make(map[string]any)

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		if idx := index(jsonTag, ','); idx != -1 {
			jsonTag = jsonTag[:idx]
		}

		result[jsonTag] = value.Interface()
	}

	return result, nil
}

// Helper function to find index of rune in string
func index(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	spltd := strings.Split(email, "@")
	if len(spltd[0]) < 2 {
		return fmt.Sprintf("%s...@%s", spltd[0][:1], spltd[1])
	}
	return fmt.Sprintf("%s...@%s", spltd[0][:2], spltd[1])
}

func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}
	return fmt.Sprintf("%s....%s", phone[:6], phone[len(phone)-1:])
}

func GetStringIfNotNull(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func StringSetNullIfEmpty(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
