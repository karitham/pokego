{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  outputs = {nixpkgs, ...}: let
    systems = [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "armv6l-linux"
      "armv7l-linux"
    ];
    forAllSystems = f: nixpkgs.lib.genAttrs systems f;
  in {
    packages = forAllSystems (system: rec {
      pokego = let
        pkgs = import nixpkgs {inherit system;};
      in
        pkgs.buildGoModule {
          pname = "pokego";
          version = "devel";
          src = ./.;
          vendorHash = "sha256-Eykg/qGqWA+qxeFPAhd0BERHtLj5X7kMQo/IPp1yRU4=";
          env.CGO_ENABLED = 0;
          flags = ["-trimpath"];
          ldflags = [
            "-s"
            "-w"
            "-extldflags -static"
          ];
          meta = {
            description = "Command-line tool that lets you display Pokémon sprites in color directly in your terminal.";
            homepage = "https://github.com/karitham/pokego";
            mainProgram = "pokego";
          };
        };
      default = pokego;
    });
    devShells = forAllSystems (system: {
      default = let
        pkgs = import nixpkgs {inherit system;};
      in
        pkgs.mkShell {
          buildInputs = [
            pkgs.go
          ];
        };
    });
  };
}
