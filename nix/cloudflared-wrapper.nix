{
  pkgs,
  targets ? [
    "x86_64-linux"
    "aarch64-linux"
    "x86_64-darwin"
    "aarch64-darwin"
    "x86_64-windows"
  ],
}:
let
  targetMeta =
    system:
    if system == "x86_64-linux" then
      {
        goos = "linux";
        goarch = "amd64";
        ext = "";
      }
    else if system == "aarch64-linux" then
      {
        goos = "linux";
        goarch = "arm64";
        ext = "";
      }
    else if system == "x86_64-darwin" then
      {
        goos = "darwin";
        goarch = "amd64";
        ext = "";
      }
    else if system == "aarch64-darwin" then
      {
        goos = "darwin";
        goarch = "arm64";
        ext = "";
      }
    else if system == "x86_64-windows" then
      {
        goos = "windows";
        goarch = "amd64";
        ext = ".exe";
      }
    else
      throw "Unsupported target system: ${system}";
in
flags:
let
  meta = map targetMeta targets;
  buildCall = pkgs.lib.strings.concatMapStringsSep "\n" (
    x: "build ${x.goos} ${x.goarch} ${x.ext}"
  ) meta;
in
pkgs.runCommand "cloudflared-wrapper"
  {
    CGO_ENABLED = "0";
  }
  #sh
  ''
    mkdir -p $out
    export HOME=/tmp

    mkdir cloudflared-wrapper
    cd cloudflared-wrapper
    cp ${../cloudflared-wrapper/go.mod} go.mod
    cp ${../cloudflared-wrapper/main.go} main.go

    build() {
      local goos="$1"
      local goarch="$2"
      local ext="$3"

      local outFile="$out/bin/cloudflared-wrapper-$goos-$goarch$ext"

      GOOS="$goos" GOARCH="$goarch" ${pkgs.lib.getExe pkgs.go} build -trimpath -ldflags='-s -w -X "main.flags=${flags}"' -o "$outFile" .
    }

    ${buildCall}
  ''
