# Decisões de arquitetura (recomendado)

Access Token (JWT): curto (ex: 15 min), usado em toda requisição.
Refresh Token: longo (ex: 7 dias), não precisa ser JWT (pode ser string aleatória/criptograficamente segura).
Revogação & sessões: refresh token persistido (Redis ou Postgres) + rotação (refresh troca por outro).
JWT assinado: HS256 (segredo) ou RS256/EdDSA (chave privada).
Claims mínimas: sub (userID), sid (sessionID), iat, exp, iss, aud, jti (opcional).
Em DDD, “autenticação” costuma viver num bounded context de Identity/Access. Mesmo que você coloque no monólito, mantenha o módulo separado.

Casos de uso

1) Login
Fluxo:
Buscar usuário por email.
Validar senha (hash).
Checar status (bloqueado etc).
Checar limite de sessões (opcional).
Criar Session com refresh token aleatório, salvar hash.
Gerar access JWT com sub=userID e sid=sessionID.
Retornar tokens.
DTO retorno:
accessToken, accessExpiresAt
refreshToken, refreshExpiresAt
sessionId

2) Refresh (rotação)
Fluxo:
Receber sessionId + refreshToken.
Buscar sessão.
Validar: não revogada, não expirada.
Comparar hash do refresh recebido.
Gerar novo refresh (rotacionar) + atualizar hash/expiração.
Gerar novo access JWT.
Retornar novos tokens.
Se refresh falhar, é comum revogar a sessão (para segurança) dependendo do seu modelo.

3) Logout
Revoke(sessionId).

4) Logout em todos dispositivos
RevokeAllByUser(userID).

5) ValidateAccess (para middleware/interceptor)
VerifyAccessToken
checar se sid ainda está ativo (opcional, depende do custo/latência)
alternativa: checar no Redis rapidamente.
Infra: JWT Service (implementação)
Boas práticas:
iss fixo (ex: actajus-identity)
aud fixo (ex: actajus-web)
exp curto
kid se rotacionar chaves
alg travado (não aceitar none)
Claims mínimas recomendadas:
sub (userID)
sid (sessionID)
exp, iat
iss, aud
roles (se você já tiver RBAC)
Infra: Refresh token seguro
Gerar com crypto/rand (32+ bytes) e base64url.
Salvar somente hash (ex: SHA-256 do token).
Comparar com hash constante.
Assim, se o Redis/vault vazar, atacante não usa refresh.
HTTP (exemplo de endpoints)
POST /auth/login
POST /auth/refresh
POST /auth/logout
POST /auth/logout-all
Cookies vs JSON
Para web browser: refresh em cookie HttpOnly Secure SameSite.
Para mobile/API: pode retornar refresh no body e armazenar seguro no client.
Access token geralmente em Authorization: Bearer ....
Middleware / gRPC interceptor
Extrair Authorization
VerifyAccessToken
Injetar userID, sessionID, roles no context.Context
(opcional) checar sessão ativa no Redis: GET session:sid -> ok

Banco/Redis (modelagem)
Redis (rápido para sessão)
session:{sid} -> {userId, refreshHash, exp, revoked=false}
índice por usuário (set):
user_sessions:{userId} -> {sid1, sid2...}
Postgres (auditoria/relatórios)
sessions table com colunas citadas acima

Você pode usar:
Redis como cache de sessão + Postgres como fonte de verdade
ou só Redis (mais simples) se você aceita a natureza in-memory.
Segurança (checklist prático)
Rate limit no login e refresh (por IP e por user/email).
Lockout progressivo (opcional).
Refresh token rotação obrigatória.
Invalidar sessão ao suspeitar de reuse de refresh (refresh reuse detection).
TLS obrigatório.
Logs com cuidado: nunca logar tokens.
Se for cookie: HttpOnly, Secure, SameSite=Lax/Strict, path do refresh restrito (/auth/refresh).
Onde cada regra fica (Clean / DDD)
Regras de negócio (limite de sessões, user bloqueado, sessão expirada/revogada): domain + orquestração no usecase.
JWT/Hash/Redis/Postgres: infrastructure.
Handlers: só traduzem request/response e chamam usecase.
Próximo passo (pra eu te entregar “com código”)
Se você quiser, eu posso te mandar um exemplo completo com:
entidades User/Session, portas, usecases Login/Refresh/Logout,
implementação JWT (lib padrão ou github.com/golang-jwt/jwt/v5),
bcrypt,
repo Redis (ou Postgres com pgxpool),
handler HTTP (chi) e interceptor gRPC (connect-go).
Só me diga qual stack você prefere:
HTTP (chi) ou gRPC (connect-go) (ou ambos)
sessão em Redis ou Postgres (ou Redis + Postgres)
assinar JWT com HS256 (segredo) ou RS256/EdDSA (chave).
