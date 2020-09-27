namespace go user_service

struct LoginRequest {
1: string username;
2: string password;
}

struct LoginResponse {
1: string msg;
}

service User {
LoginResponse checkPassword(1: LoginRequest req);
}
