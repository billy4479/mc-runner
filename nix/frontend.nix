{
  stdenv,
  pnpm_11,
  pnpmConfigHook,
  fetchPnpmDeps,
  nodejs,
  ...
}:
stdenv.mkDerivation (finalAttrs: {
  pname = "frontend";
  version = "1.0.0";
  src = ../frontend;
  nativeBuildInputs = [
    nodejs
    pnpm_11
    pnpmConfigHook
  ];

  pnpmDeps = fetchPnpmDeps {
    inherit (finalAttrs) pname version src;
    hash = "sha256-1vBtJfN5MnMbvpEzm19BIFU5jkSg8GaE7y+riItc+uM=";
    fetcherVersion = 3;
  };

  buildPhase = # sh
    ''
      pnpm build
    '';
  installPhase = # sh
    ''
      mv dist $out
    '';
})
