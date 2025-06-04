const API_ENDPOINT = import.meta.env.DEV ? "http://localhost:4479/api" : "/api";

export async function ping() {
  await fetch(`${API_ENDPOINT}/ping`)
    .then((res) => res.text())
    .then(console.log);
}
