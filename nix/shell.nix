pkgs:
pkgs.mkShell {
  packages = with pkgs; [
    cloudflared
    podman
    turso-cli
    air

    gopls
    golangci-lint
    (go-migrate.overrideAttrs (oldAttrs: {
      tags = [ "sqlite3" ];
    }))
    delve

    typescript-language-server
    svelte-language-server
    tailwindcss-language-server

    nixd
  ];

  nativeBuildInputs = with pkgs; [
    go
    sqlc

    nodejs_latest
    nodePackages.pnpm

    sqlitebrowser
  ];
}
