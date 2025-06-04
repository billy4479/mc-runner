{
  stdenv,
  pnpm,
  nodejs,
  ...
}:
stdenv.mkDerivation (finalAttrs: {
  pname = "frontend";
  version = "1.0.0";
  src = ../.;
  nativeBuildInputs = [
    nodejs
    pnpm.configHook
  ];

  pnpmDeps = pnpm.fetchDeps {
    inherit (finalAttrs) pname version src;
    hash = "sha256-GMbYQa3AsLsCzeQOfyuAQOcW0o9GgEArfCvvowId+Wo=";
    sourceRoot = "${finalAttrs.src}/frontend";
  };

  pnpmRoot = "frontend";

  buildPhase = # sh
    ''
      pnpm -C frontend build
    '';
  installPhase = # sh
    ''
      mv frontend/dist $out
    '';
})
