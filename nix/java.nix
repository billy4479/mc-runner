{
  stdenvNoCC,
  lib,

  jre8_headless,
  jdk17_headless,
  jdk21_headless,
  ...
}:
let
  link =
    lib.strings.concatMapAttrsStringSep "\n" (n: v: "ln -s ${lib.getExe' v "java"} $out/bin/java${n}")
      {
        "8" = jre8_headless;
        "17" = jdk17_headless;
        "21" = jdk21_headless;
      };
in
stdenvNoCC.mkDerivation {
  name = "mc-java";

  dontBuild = true;
  dontUnpack = true;

  installPhase = # sh
    ''
      mkdir -p $out/bin
      ${link}
    '';
}
