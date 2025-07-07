export async function ping() {
  await fetch("/api/ping")
    .then((res) => res.text())
    .then(console.log);
}

export async function register(token: string, name: string) {
  await fetch("/api/auth/register", {
    method: "POST",
    credentials: "include",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify({ invitation_token: token, name: name }),
  })
    .then((res) => res.json())
    .then(console.log);
}
