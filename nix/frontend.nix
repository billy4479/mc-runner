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
    hash = "sha256-1Gu5eLh5JaXVPnimkmQhZkbxFTg7jtJymD7c3Y0BlmE=";
    fetcherVersion = 4;
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
