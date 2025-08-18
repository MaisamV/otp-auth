# Project Summary
- An OTP-based authentication backend service using Golang 1.24.4

## Flows
### 1. Send OTP Flow
- User requests to send code API using phone number with format +<country_code><phone_number> or 0<phone_number> (the <phone_number> MUST be a 10 length string which starts with 9) (example of correct number 09123456789).
- System receives the request and extract session_id either from cookies or JSON request, create a PhoneNumber VO or Entity using the given phone number which has validate method and validates it, if the phone number is valid check rate limiting by incrementing a counter in Redis with key "rate_limit:<phone_number>" with 10 minutes TTL and verify it's not more than 3, then generate a random OTP, send the OTP with it's related phone number to OtpSender interface, hash the OTP for security reasons (users with access to redis can login as someone else), if user sent no session_id create a random text as session_id and store the <session_id>-<hashed OTP> in redis with key <phone_number> as key with a 2 minute TTL and return session_id both in json response and also set it as httpOnly cookie without expiration.
- **NOTE:**
  Expire hashed OTP after 2 minutes (configurable in a yaml config file).
  Limit OTP requests: Max 3 requests per phone number within 10 minutes using Redis counter.
  session_id acts as an anti-brute force mechanism - attackers cannot attempt login without first going through the send OTP flow.

### 2. Login/Register Flow
- User requests to login/register using session_id (which can be sent in Json or read from cookies) ,phone number and OTP; system hashes OTP, checks the redis with key <phone_number> and match the value with <session_id>-<hashed OTP> and If it matched and not expired then create the user in postgres if not already exists (or ignore this step if user already exists) and create a JWT token as access token using JwtService interface with claims structure and a long opaque string as refresh token; save refresh token hash, session_id, user_id, created_at, expires_at, last_used, revoked and revoke_reason in postgres and send these tokens in API json response and also set these as httpOnly cookies too with CSRF protection.
**NOTE:**
  Access token expires in 15 minutes, refresh token expires in 30 days (configurable)
  If session_id is missing, return error 400 "session_id required - please call send OTP first"
  Created users has scope, empty scope means normal users, and superadmin scope has access to everything
  OtpSender implementation should just print the OTP and phone number to console
  JwtService interface supports 2 methods, one that receive the claims and returns JWT token string (and an error if any might happen) and one that verifies JWT tokens.
  JwtService implementation constructor starts to read ecdsa private and public keys from keys folder, if even one of them is missing, then it generates ecdsa private and public key in keys folder. it uses these keys to generate JWT token and verify issued tokens.
  JWT Claims structure: {"sub": userID, "client_id": clientID, "scopes": scopes, "iat": now.Unix(), "exp": expiry.Unix(), "iss": issuer, "jti": tokenID}
  Refresh tokens table includes user_id to enable bulk invalidation of all user's tokens when needed
  Create a Hash function and random string generator for refresh token and define them in common or utils to prevent duplicate code

### 3. Refresh Flow
- User requests to refresh in order to renew access token and refresh_token, API should receive refresh_token and session_id (which can be sent in Json or read from cookies); system finds database record where the refresh_token_hash matches the hash of the provided refresh_token AND the session_id matches the provided session_id AND not revoked yet, extracts user_id from that record, then updates it with revoked=true and last_used=NOW() and revoke_reason='REFRESH', if any row updated, then the record existed so create a new refresh token and access token and create a new record with hash of refresh token, session_id, and the extracted user_id and return these tokens to user both in JSON response and set as cookies just like login/register response.


### 4. Logout Flow
- User requests to logout using refresh_token (which can be sent in Json or read from cookies); system updates database record which has the same refresh token hash and not revoked yet with revoked=true and last_used=NOW() and revoke_reason='LOGOUT', clear httpOnly cookies and return success response.

### 5. Admin User Management Flow
- Retrieve single user details.
- Retrieve list of users with Pagination and Search (by phone number and registration date fields).

## Error Handling

### Send OTP Flow Errors:
- **400 Bad Request**: Invalid phone number format
- **429 Too Many Requests**: Rate limit exceeded (more than 3 requests in 10 minutes)
- **500 Internal Server Error**: OTP generation or Redis storage failure

### Login/Register Flow Errors:
- **400 Bad Request**: Missing session_id, phone number, or OTP
- **401 Unauthorized**: Invalid OTP or session_id mismatch
- **404 Not Found**: No OTP found for phone number (expired or never sent)
- **410 Gone**: OTP expired
- **500 Internal Server Error**: Database or JWT generation failure

### Refresh Flow Errors:
- **400 Bad Request**: Missing refresh token or session_id
- **401 Unauthorized**: Invalid or expired refresh token, or session_id mismatch
- **403 Forbidden**: Refresh token already revoked
- **500 Internal Server Error**: Database or JWT generation failure

### Logout Flow Errors:
- **400 Bad Request**: Missing refresh token
- **401 Unauthorized**: Invalid refresh token
- **500 Internal Server Error**: Database update failure

## Architecture
- Use a clean and hexagonal architecture as main architecture with SOLID principles in mind.
- use separate reader, writer interfaces for repositories.
- Anywhere that you detect duplicate code, you should create a function or method of that code and move it to the part of project that is responsible for that functionality. 
- separate domain, application, infrastructure layers based on Clean architecture
- Ensure clear separation of responsibilities in code.

## CI/CD
- The project must have a dockerfile
- Set up the project with docker-compose.

## Database
- use Redis to keep hashed OTPs of users and rate limits.
- Starting with postgres while the load is not too much and when the users and load increases we can add citus for user and refresh table management.

## Documentation
- Document all REST APIs with OpenAPI.
- README.md should contain all these titles: Project summary, How to run with Docker compose, How to run locally.
- WHY.md should describe why we chose this tech stack, architecture, databases, encryption method ecdsa, session_id adding to OTP flow, etc and also compare it with alternative ones.