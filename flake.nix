{
  description = "greed: go rss feed reader";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      # Provide some binary packages for selected system types.
      packages.${system} = {
        greed = pkgs.buildGoModule {
          pname = "greed";
          version = "1.0.0";

          src = ./.;

          vendorHash = "sha256-0RHgeJZCoRCwImYHrnpaJ9pdrZnnRpnNZu5HqLPKJb4=";
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gotools
          go-tools
        ];
      };

      defaultPackage.${system} = self.packages.${system}.greed;
    };
}
