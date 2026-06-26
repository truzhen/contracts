package secrets

// SecretRef 用于引用 SecureStore 中的密钥，绝不包含明文
type SecretRef struct {
	SecretID string `json:"secret_id"`
}

// SensitivePayload 代表经过加密或脱敏的负载
type SensitivePayload struct {
	EncryptedData []byte    `json:"encrypted_data"`
	SecretRef     SecretRef `json:"secret_ref"`
}
