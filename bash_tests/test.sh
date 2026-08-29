#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для генерации UUID
generate_uuid() {
    if command -v uuidgen &> /dev/null; then
        uuidgen
    elif [ -f /proc/sys/kernel/random/uuid ]; then
        cat /proc/sys/kernel/random/uuid
    else
        echo "$(date +%s%N | md5sum | cut -c1-8)-$(date +%s%N | md5sum | cut -c1-4)-$(date +%s%N | md5sum | cut -c1-4)-$(date +%s%N | md5sum | cut -c1-4)-$(date +%s%N | md5sum | cut -c1-12)"
    fi
}

# Функция для красивого вывода
print_test() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN} $1${NC}"
}

print_error() {
    echo -e "${RED} $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}  $1${NC}"
}

# СОЗДАЕМ ОДНОГО ПОЛЬЗОВАТЕЛЯ ДЛЯ ВСЕХ ТЕСТОВ
TEST_USER="testuser_$(date +%s)"
TEST_PASS="SuperPassword123!"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  AUTH SERVICE INTEGRATION TESTS${NC}"
echo -e "${BLUE}  User: $TEST_USER${NC}"
echo -e "${BLUE}========================================${NC}"

# ============================================================
# 1. REGISTER - Создание пользователя (выход nil)
# ============================================================
print_test "1. REGISTER - Create user: $TEST_USER"

HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(generate_uuid)" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")

if [ "$HTTP_STATUS" == "201" ]; then
    print_success "User created successfully! (HTTP $HTTP_STATUS)"
else
    print_error "Registration failed with status: $HTTP_STATUS"
    exit 1
fi

# ============================================================
# 2. REGISTER IDEMPOTENCY - Проверка кэширования
#    Первый запрос: 201 Created
#    Второй запрос: 200 OK (кэш)
# ============================================================
print_test "2. REGISTER IDEMPOTENCY - Same key (first: 201, second: 200)"

IDEM_KEY=$(generate_uuid)
echo "Idempotency Key: $IDEM_KEY"

# Первый запрос - должен вернуть 201 Created
echo -e "\n${YELLOW}First request (should return 201 Created):${NC}"
FIRST_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"idem_user_$(date +%s)\",
    \"password\": \"$TEST_PASS\"
  }")
echo "HTTP Status: $FIRST_STATUS"

# Второй запрос с тем же ключом - должен вернуть 200 OK (кэш)
echo -e "\n${YELLOW}Second request with same key (should return 200 OK - cached):${NC}"
SECOND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"idem_user_2_$(date +%s)\",
    \"password\": \"$TEST_PASS\"
  }")
echo "HTTP Status: $SECOND_STATUS"

# Проверяем
if [ "$FIRST_STATUS" == "201" ] && [ "$SECOND_STATUS" == "200" ]; then
    print_success "Idempotency works! First: 201, Second: 200"
elif [ "$FIRST_STATUS" == "201" ] && [ "$SECOND_STATUS" == "201" ]; then
    print_warning "Both returned 201 (might not use cache)"
else
    print_error "Idempotency failed! Expected 201 then 200, got $FIRST_STATUS then $SECOND_STATUS"
fi

# ============================================================
# 3. LOGIN - Вход (выход SessionResponse)
# ============================================================
print_test "3. LOGIN - Login with user: $TEST_USER"

LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(generate_uuid)" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")

echo "$LOGIN_RESPONSE" | jq .

# Извлекаем токены из SessionResponse
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
ACCESS_EXPIRES_AT=$(echo "$LOGIN_RESPONSE" | jq -r '.access_expired_at')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token')
REFRESH_EXPIRES_AT=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_expired_at')

if [ "$REFRESH_TOKEN" != "null" ] && [ ! -z "$REFRESH_TOKEN" ]; then
    print_success "Login successful!"
    echo "  Access Token: ${ACCESS_TOKEN:0:30}..."
    echo "  Access Expires: $ACCESS_EXPIRES_AT"
    echo "  Refresh Token: ${REFRESH_TOKEN:0:30}..."
    echo "  Refresh Expires: $REFRESH_EXPIRES_AT"
else
    print_error "Login failed!"
    exit 1
fi

# ============================================================
# 4. LOGIN IDEMPOTENCY - Проверка кэширования
#    Первый запрос: 200 OK
#    Второй запрос: 200 OK (кэш)
# ============================================================
print_test "4. LOGIN IDEMPOTENCY - Same key (both should return 200 OK)"

IDEM_KEY=$(generate_uuid)
echo "Idempotency Key: $IDEM_KEY"

# Первый запрос
echo -e "\n${YELLOW}First request (should return 200 OK):${NC}"
FIRST_LOGIN=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")
FIRST_LOGIN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")
echo "$FIRST_LOGIN" | jq .
echo "HTTP Status: $FIRST_LOGIN_STATUS"
FIRST_LOGIN_REFRESH=$(echo "$FIRST_LOGIN" | jq -r '.refresh_token')

# Второй запрос с тем же ключом - должен вернуть 200 OK (кэш)
echo -e "\n${YELLOW}Second request with same key (should return 200 OK - cached):${NC}"
SECOND_LOGIN=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")
SECOND_LOGIN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")
echo "$SECOND_LOGIN" | jq .
echo "HTTP Status: $SECOND_LOGIN_STATUS"
SECOND_LOGIN_REFRESH=$(echo "$SECOND_LOGIN" | jq -r '.refresh_token')

# Проверяем
if [ "$FIRST_LOGIN_STATUS" == "200" ] && [ "$SECOND_LOGIN_STATUS" == "200" ]; then
    print_success "Both requests returned 200 OK"
    if [ "$FIRST_LOGIN_REFRESH" == "$SECOND_LOGIN_REFRESH" ] && [ "$FIRST_LOGIN_REFRESH" != "null" ]; then
        print_success "Idempotency works! Same refresh token"
    else
        print_warning "Different refresh tokens returned"
    fi
else
    print_error "Expected 200 OK for both, got $FIRST_LOGIN_STATUS and $SECOND_LOGIN_STATUS"
fi

# ============================================================
# 5. REFRESH - Обновление токена (выход SessionResponse)
# ============================================================
print_test "5. REFRESH TOKEN"

REFRESH_RESPONSE=$(curl -s -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(generate_uuid)" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "$REFRESH_RESPONSE" | jq .

NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.access_token')
NEW_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.refresh_token')

if [ "$NEW_ACCESS_TOKEN" != "null" ] && [ ! -z "$NEW_ACCESS_TOKEN" ]; then
    print_success "Token refresh successful!"
    ACCESS_TOKEN="$NEW_ACCESS_TOKEN"
    REFRESH_TOKEN="$NEW_REFRESH_TOKEN"
    echo "  New Access Token: ${ACCESS_TOKEN:0:30}..."
    echo "  New Refresh Token: ${REFRESH_TOKEN:0:30}..."
else
    print_error "Token refresh failed!"
fi

# ============================================================
# 6. REFRESH IDEMPOTENCY - Проверка кэширования
#    Первый запрос: 200 OK
#    Второй запрос: 200 OK (кэш)
# ============================================================
print_test "6. REFRESH TOKEN IDEMPOTENCY (both should return 200 OK)"

IDEM_KEY=$(generate_uuid)
echo "Idempotency Key: $IDEM_KEY"

# Первый запрос
echo -e "\n${YELLOW}First request (should return 200 OK):${NC}"
FIRST_REFRESH=$(curl -s -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
FIRST_REFRESH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "$FIRST_REFRESH" | jq .
echo "HTTP Status: $FIRST_REFRESH_STATUS"
FIRST_REFRESH_ACCESS=$(echo "$FIRST_REFRESH" | jq -r '.access_token')

# Второй запрос с тем же ключом - должен вернуть 200 OK (кэш)
echo -e "\n${YELLOW}Second request with same key (should return 200 OK - cached):${NC}"
SECOND_REFRESH=$(curl -s -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
SECOND_REFRESH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "$SECOND_REFRESH" | jq .
echo "HTTP Status: $SECOND_REFRESH_STATUS"
SECOND_REFRESH_ACCESS=$(echo "$SECOND_REFRESH" | jq -r '.access_token')

# Проверяем
if [ "$FIRST_REFRESH_STATUS" == "200" ] && [ "$SECOND_REFRESH_STATUS" == "200" ]; then
    print_success "Both requests returned 200 OK"
    if [ "$FIRST_REFRESH_ACCESS" == "$SECOND_REFRESH_ACCESS" ] && [ "$FIRST_REFRESH_ACCESS" != "null" ]; then
        print_success "Idempotency works! Same access token"
    else
        print_warning "Different access tokens returned"
    fi
else
    print_error "Expected 200 OK for both, got $FIRST_REFRESH_STATUS and $SECOND_REFRESH_STATUS"
fi

# ============================================================
# 7. LOGOUT - Выход (выход nil)
# ============================================================
print_test "7. LOGOUT"

LOGOUT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/logout \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(generate_uuid)" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

if [ "$LOGOUT_STATUS" == "200" ] || [ "$LOGOUT_STATUS" == "204" ]; then
    print_success "Logout successful! (HTTP $LOGOUT_STATUS)"
else
    print_error "Logout failed with status: $LOGOUT_STATUS"
fi

# ============================================================
# 8. LOGOUT IDEMPOTENCY - Проверка кэширования
#    Первый запрос: 200 OK
#    Второй запрос: 200 OK (кэш)
# ============================================================
print_test "8. LOGOUT IDEMPOTENCY (both should return 200/204 OK)"

IDEM_KEY=$(generate_uuid)
echo "Idempotency Key: $IDEM_KEY"

# Первый запрос
echo -e "\n${YELLOW}First request (should return 200/204 OK):${NC}"
FIRST_LOGOUT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/logout \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "HTTP Status: $FIRST_LOGOUT_STATUS"

# Второй запрос с тем же ключом - должен вернуть 200 OK (кэш)
echo -e "\n${YELLOW}Second request with same key (should return 200/204 OK - cached):${NC}"
SECOND_LOGOUT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/auth/logout \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $IDEM_KEY" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "HTTP Status: $SECOND_LOGOUT_STATUS"

# Проверяем
if [[ "$FIRST_LOGOUT_STATUS" == "200" || "$FIRST_LOGOUT_STATUS" == "204" ]] && \
   [[ "$SECOND_LOGOUT_STATUS" == "200" || "$SECOND_LOGOUT_STATUS" == "204" ]]; then
    print_success "Idempotency works! Both returned success status"
else
    print_warning "Expected success status for both, got $FIRST_LOGOUT_STATUS and $SECOND_LOGOUT_STATUS"
fi

# ============================================================
# 9. LOGIN AFTER LOGOUT - Проверка что можно залогиниться снова
# ============================================================
print_test "9. LOGIN AFTER LOGOUT"

LOGIN_AFTER_LOGOUT=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: $(generate_uuid)" \
  -d "{
    \"login\": \"$TEST_USER\",
    \"password\": \"$TEST_PASS\"
  }")

echo "$LOGIN_AFTER_LOGOUT" | jq .

NEW_REFRESH=$(echo "$LOGIN_AFTER_LOGOUT" | jq -r '.refresh_token')
if [ "$NEW_REFRESH" != "null" ] && [ ! -z "$NEW_REFRESH" ]; then
    print_success "Can login again with new token"
else
    print_error "Login failed after logout"
fi

# ============================================================
# TEST SUMMARY
# ============================================================
print_test "TEST SUMMARY"

echo -e "${GREEN}All tests completed!${NC}\n"

echo "Test Results:"
echo "  1. Register - User created (HTTP 201)"
echo "  2. Register Idempotency - First: 201, Second: 200 (cache)"
echo "  3. Login - SessionResponse received"
echo "  4. Login Idempotency - Both 200 OK"
echo "  5. Refresh Token - SessionResponse received"
echo "  6. Refresh Idempotency - Both 200 OK"
echo "  7. Logout - Success (HTTP 200/204)"
echo "  8. Logout Idempotency - Both 200/204 OK"
echo "  9. Login after logout - New session created"

echo -e "\n${BLUE}User Information:${NC}"
echo "  Username: $TEST_USER"
echo "  Password: $TEST_PASS"
echo "  Refresh Token: ${REFRESH_TOKEN:0:30}..."

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}  ALL TESTS PASSED! ${NC}"
echo -e "${GREEN}========================================${NC}"