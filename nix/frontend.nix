{
  stdenv,
  pnpm,
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
    pnpm
    pnpmConfigHook
  ];

  pnpmDeps = fetchPnpmDeps {
    inherit (finalAttrs) pname version src;
    hash = "sha256-7NWM3RYGiPZC87zbykCoWStSlmAdn1S78q50tvQbifI=";
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
