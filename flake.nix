{
  description = "Example Go development environment for Zero to Nix";

  # Flake inputs
  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2405.*.tar.gz";
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      # Systems supported
      allSystems = [
        "x86_64-linux"   # 64-bit Intel/AMD Linux
        "aarch64-linux"  # 64-bit ARM Linux
        "x86_64-darwin"  # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS (Apple Silicon)
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            go_1_22   # The Go CLI
            gotools   # Go tools like goimports, godoc, and others
            sqlc      # SQL code generator

            # OS-specific tools
          ] ++ (if pkgs.stdenv.isLinux then [ inotify-tools ] else [])
            ++ (if pkgs.stdenv.isDarwin then [ fswatch ] else []);
        };
      });
    };
}
