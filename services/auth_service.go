package services

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash from a plain text password.
// The resulting hash can be safely stored in the database.
func HashPassword(password string) (string, error) {

	// Generate a bcrypt hash using the default cost factor.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	// Convert the hash from bytes to string before returning it.
	return string(hashedPassword), nil
}

// CheckPassword compares a plain text password with a stored bcrypt hash.
// It returns true if the password matches the hash.
func CheckPassword(password string, hashedPassword string) bool {

	// Compare the provided password against the stored hash.
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	// Return true only when the comparison succeeds.
	return err == nil
}