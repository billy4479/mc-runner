pkgs:
pkgs.mkShell {
  packages = with pkgs; [
    cloudflared
    podman
    turso-cli
    air

    gopls
    golangci-lint

    typescript-language-server
    svelte-language-server
    tailwindcss-language-server

    nixd
  ];

  nativeBuildInputs = with pkgs; [
    go

    nodejs_latest
    nodePackages.pnpm
  ];
}
