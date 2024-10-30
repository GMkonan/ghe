{
  description = "ghe dev env golang";

  inputs =
    {
      nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    };

  outputs = { self, nixpkgs, ... }@inputs:
    let
      system = "aarch64-darwin";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.aarch64-darwin.default =
        pkgs.mkShell
          {
            nativeBuildInputs = with pkgs; [
             go 
            ];
          shellHook = ''
            echo "Started your ghe dev env shell"
            '';
          };
    };
}
