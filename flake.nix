{
  description = "Go + cgo + X11 dev shell";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    {
      devShells.x86_64-linux.default =
        let
          pkgs = nixpkgs.legacyPackages.x86_64-linux;
        in
        pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.libx11
            pkgs.pkg-config
            pkgs.wl-clipboard
            pkgs.xclip
            pkgs.xsel
          ];
        };
    };
}
