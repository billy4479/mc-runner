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

    (sops.overrideAttrs (
      finalAttrs: previousAttrs: {
        patches = [
          (fetchpatch {
            url = "https://github.com/getsops/sops/pull/1914.diff";
            hash = "sha256-80L8JbzrZHhH8608B1YO+rlpYIYWD6gIvfqjkNATmoA=";
          })
        ];
      }
    ))
  ];

  nativeBuildInputs = with pkgs; [
    go
    sqlc

    nodejs_latest
    nodePackages.pnpm

    sqlitebrowser
  ];
}
