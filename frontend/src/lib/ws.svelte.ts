export class ServerSocket {
  private _ws: WebSocket | null = null;
  #isConnected = $state(false);
  #serverVersionString = $state("");

  private _reconnectTimeout: number | null = null;

  get isConnected() {
    return this.#isConnected;
  }
  get serverVersionString() {
    return this.#serverVersionString;
  }

  constructor() {}

  close() {
    if (this._ws) {
      this._ws.close(1000);
      if (this._reconnectTimeout) clearTimeout(this._reconnectTimeout);

      console.log("closing cleanely from client");
    }
  }

  connect() {
    if (this._ws !== null) return;

    this._ws = new WebSocket("/api/ws");
    const ws = this._ws;
    console.log("connected ws");

    ws.addEventListener("close", (ev) => {
      this.#isConnected = false;
      this._ws = null;

      if (ev.code != 1000) {
        console.warn("WS closed", ev);
        this._reconnectTimeout = setTimeout(() => {
          console.log("attempting reconnection");
          this.connect();
        }, 1000);
      }
    });

    ws.addEventListener("error", console.error);
    ws.addEventListener("message", (message) => {
      const msg = JSON.parse(message.data);
      console.log(msg);

      switch (msg.type) {
        case "version":
          this.#serverVersionString = msg.data;
          break;
        default:
          console.warn("unknown message type");
          break;
      }
    });
    ws.addEventListener("open", () => {
      ws.send("ping");
      this.#isConnected = true;
    });
  }
}
