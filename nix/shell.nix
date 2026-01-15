pkgs:
pkgs.mkShell {
  packages = with pkgs; [
    podman
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

    prettier

    nixd

    sops
  ];

  buildInputs = with pkgs; [
    mc-java
  ];

  nativeBuildInputs = with pkgs; [
    go
    sqlc

    nodejs_latest
    nodePackages.pnpm

    sqlitebrowser
  ];
}
