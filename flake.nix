{
  description = "ProblemService gRPC development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            buf
            grpcurl

            go
            go-tools
            gotest
          ];

          shellHook = ''
            echo "================================================================"
            echo " command:"
            echo "   - buf           : Compile .proto files and generate Go code"
            echo "   - go            : Run Go commands (build, test, etc.)"
            echo "   - gotest        : Run Go tests in the current directory"
            echo "   - grpcurl       : Test gRPC services from the command line"
            echo "================================================================"
          '';
        };
      }
    );
}
