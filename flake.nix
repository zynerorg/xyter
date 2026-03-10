{
  description = "Nix flake for the XYter project with Go + frontend build";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        # Define apps first to avoid self-reference issues
        appsDef = rec {
          xyter-api = {
            type = "app";
            program = "${self.packages.${system}.xyter-api}/bin/xyter-api";
          };
          xyter-cli = {
            type = "app";
            program = "${self.packages.${system}.xyter-cli}/bin/xyter-cli";
          };
          xyter-bot = {
            type = "app";
            program = "${self.packages.${system}.xyter-bot}/bin/xyter-bot";
          };
          xyter-all = {
            type = "app";
            program = "${self.packages.${system}.xyter-all}/bin";
          };
        };
      in
      rec {
        #############################################################################
        # Go binaries
        #############################################################################
        packages.xyter-api = pkgs.buildGoModule {
          pname = "xyter-api";
          version = "0.1.0";
          src = ./.;
          #goPackagePath = "github.com/yourusername/xyter";
          subPackages = [ "./cmd/xyter-api" ];
          vendorHash = "sha256-/FvLF0BB6jl4ZsdtgWr/PlcvHukcR5vKKJfzQ5+M6f0=";

          installPhase = ''
            mkdir -p $out/bin
            find $src
            cp result/bin/xyter-api $out/bin/
          '';
        };

        packages.xyter-cli = pkgs.buildGoModule {
          pname = "xyter-cli";
          version = "0.1.0";
          src = ./.;
          #goPackagePath = "github.com/yourusername/xyter";
          subPackages = [ "./cmd/xyter-cli" ];
          vendorHash = "sha256-/FvLF0BB6jl4ZsdtgWr/PlcvHukcR5vKKJfzQ5+M6f0=";

          installPhase = ''
            mkdir -p $out/bin
            cp result/bin/xyter-cli $out/bin/
          '';
        };

        packages.xyter-bot = pkgs.buildGoModule {
          pname = "xyter-bot";
          version = "0.1.0";
          src = ./.;
          #goPackagePath = "github.com/yourusername/xyter";
          subPackages = [ "./cmd/xyter-bot" ];
          vendorHash = "sha256-/FvLF0BB6jl4ZsdtgWr/PlcvHukcR5vKKJfzQ5+M6f0=";

          installPhase = ''
            mkdir -p $out/bin
            cp result/bin/xyter-bot $out/bin/
          '';
        };

        #############################################################################
        # Frontend build (Trunk/WASM)
        #############################################################################
        packages.xyter-frontend = pkgs.stdenv.mkDerivation {
          pname = "xyter-frontend";
          version = "0.1.0";

          src = ./.;

          nativeBuildInputs = [
            pkgs.rustup
            pkgs.cargo
            pkgs.trunk
            pkgs.nodejs
          ];

          buildPhase = ''
            export CARGO_HOME=$(pwd)/tools/cargo
            export RUSTUP_HOME=$(pwd)/tools/rustup
            export PATH=$CARGO_HOME/bin:$PATH

            echo "Building frontend with trunk..."
            trunk build --release --dist ./build/static
          '';

          installPhase = ''
            mkdir -p $out/static
            cp -r ./build/static/* $out/static/
          '';
        };

        #############################################################################
        # Combined "all" package
        #############################################################################
        packages.xyter-all = pkgs.stdenv.mkDerivation {
          pname = "xyter-all";
          version = "0.1.0";

          buildInputs = [
            self.packages.${system}.xyter-api
            self.packages.${system}.xyter-cli
            self.packages.${system}.xyter-bot
            self.packages.${system}.xyter-frontend
          ];

          installPhase = ''
            mkdir -p $out/bin
            mkdir -p $out/static

            cp -r ${self.packages.${system}.xyter-api}/bin/* $out/bin/
            cp -r ${self.packages.${system}.xyter-cli}/bin/* $out/bin/
            cp -r ${self.packages.${system}.xyter-bot}/bin/* $out/bin/
            cp -r ${self.packages.${system}.xyter-frontend}/static/* $out/static/
          '';
        };

        #############################################################################
        # Apps for nix run
        #############################################################################
        apps = appsDef;
        defaultApp = apps.xyter-all;

        #############################################################################
        # Development shell
        #############################################################################
        devShell = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.rustup
            pkgs.cargo
            pkgs.trunk
            pkgs.nodejs
            pkgs.golangci-lint
            pkgs.delve
          ];

          shellHook = ''
            echo "Welcome to the XYter dev shell!"
            echo "Go, Rust, NodeJS, Trunk, and dev tools are ready"
          '';
        };
      }
    );
}
