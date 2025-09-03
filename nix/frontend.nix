{
  stdenv,
  pnpm,
  nodejs,
  ...
}:
stdenv.mkDerivation (finalAttrs: {
  pname = "frontend";
  version = "1.0.0";
  src = ../frontend;
  nativeBuildInputs = [
    nodejs
    pnpm.configHook
  ];

  pnpmDeps = pnpm.fetchDeps {
    inherit (finalAttrs) pname version src;
    hash = "sha256-8BU0Xqoz+7o+iKf7+vbRPjVnRdk5HUk1TnHbyDs1iDs=";
    fetcherVersion = 2;
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
