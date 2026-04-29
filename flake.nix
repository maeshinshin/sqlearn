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
            # Protocol Buffers コアツール
            buf # Linter, Formatter, Generatorとして強力なモダンツール

            # gRPC デバッグツール
            grpcurl # コマンドラインからgRPCリクエストを送信するツール

            # Go言語ツールチェーンとprotocプラグイン
            go
            go-tools
          ];

          shellHook = ''
            echo "================================================================"
            echo " command:"
            echo "   - buf           : Compile .proto files and generate Go code"
            echo "   - go            : Run Go commands (build, test, etc.)"
            echo "   - grpcurl       : Test gRPC services from the command line"
            echo "================================================================"
          '';
        };
      }
    );
}
