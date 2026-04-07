package jwt

import (
	"fmt"
	"testing"
)

// TestJWT 测试JWT的生成和解析
func TestJWT(t *testing.T) {
	// 测试生成token
	userID := int64(809852607525818368)
	username := "testuser"

	fmt.Println("========== 测试生成新Token ==========")
	token, err := GenToken(userID, username)
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}
	fmt.Printf("生成的Token: %s\n", token)

	// 测试解析刚生成的token
	fmt.Println("\n========== 测试解析刚生成的Token ==========")
	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("解析Token失败: %v", err)
	}
	fmt.Printf("解析成功! UserID: %d, Username: %s\n", claims.UserID, claims.Username)

	// 验证解析出的数据是否正确
	if claims.UserID != userID {
		t.Errorf("UserID不匹配: 期望 %d, 实际 %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Username不匹配: 期望 %s, 实际 %s", username, claims.Username)
	}

	fmt.Println("\n========== 测试解析旧Token ==========")
	oldToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4MDk4NTI2MDc1MjU4MTgzNjgsInVzZXJuYW1lIjoidXNlcm5hbWUiLCJleHAiOjE3NzQ5NTEyMTgsImlzcyI6ImJsdWViZWxsIn0.Ae2F9OEW5eZv3Z0NwqtBRICoZLOzZNXlxd9W9MDbbas"
	claims2, err := ParseToken(oldToken)
	if err != nil {
		fmt.Printf("解析旧Token失败: %v\n", err)
		fmt.Println("注意: 旧Token解析失败是正常的，因为密钥可能不匹配")
	} else {
		fmt.Printf("解析成功! UserID: %d, Username: %s\n", claims2.UserID, claims2.Username)
	}

	fmt.Println("\n========== 测试解析无效Token ==========")
	invalidToken := "invalid.token.string"
	_, err = ParseToken(invalidToken)
	if err != nil {
		fmt.Printf("解析无效Token返回错误(符合预期): %v\n", err)
	} else {
		t.Error("解析无效Token应该返回错误,但却成功了")
	}
}
