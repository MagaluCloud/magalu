{
  description = "Magalu Nix distribution channel";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
      mkDeps = pkgs: pkgs.stdenv.mkDerivation {
        name = "cli-deps";
        src = ./.;

        nativeBuildInputs = [ pkgs.buildPackages.go ];

        configurePhase = ''
          export GOMODCACHE=$out
          export GOPATH=$(mktemp -d)
        '';

        buildPhase = ''
          export GOPROXY=https://proxy.golang.org,direct
          go mod download
        '';

        dontInstall = true;

        outputHashMode = "recursive";
        outputHashAlgo = "sha256";
        outputHash = "M0WccPKQTg6V+7wr0kkx0E6aB7QudHtD44TqyDCRtuQ=";
      };
    in
    {
      packages = forAllSystems ({ pkgs }:
        let
          version = "0.51.0";
          deps = mkDeps pkgs;

          cli = pkgs.buildGoModule {
            pname = "mgc";
            inherit version;

            src = pkgs.lib.cleanSource ./.;

            modRoot = "./mgc/cli";

            proxyVendor = true;
            vendorHash = null;

            preBuild = ''
              export GOMODCACHE=${deps}
            '';

            tags = [ "embed" ];
            ldflags = [ "-s" "-w" "-X" "main.RawVersion=v${version}" ];
            subPackages = [ "." ];

            postInstall = ''
              mv $out/bin/cli $out/bin/mgc
            '';
          };

          # Create aliases for the same package
          packages = {
            cli = cli;
            deps = deps;
            default = cli;
          };
        in
        packages
      );

      apps = forAllSystems ({ pkgs }: {
        cli = {
          type = "app";
          program = "${self.packages.${pkgs.system}.mgc}/bin/mgc";
        };
        default = self.apps.${pkgs.system}.mgc;
      });
    };
}
