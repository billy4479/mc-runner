export async function ping() {
  await fetch("/api/ping")
    .then((res) => res.text())
    .then(console.log);
}

export async function register(token: string, name: string) {
  await fetch("/api/auth/register", {
    method: "POST",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({ invitation_token: token, name: name }),
  })
    .then((res) => res.json())
    .then(console.log);
}

export async function login(token: string) {
  await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({ login_token: token }),
  })
    .then((res) => res.json())
    .then(console.log);
}

export async function getMe() {
  return fetch("/api/auth/me").then((res) => res.json());
}

export async function logout() {
  await fetch("/api/auth/logout", {
    method: "DELETE",
  })
    .then((res) => res.json())
    .then(console.log);
}
