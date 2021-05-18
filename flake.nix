{
  description = "An open source implementation of the Tailscale coordination server";

  inputs.nixpkgs.url = github:NixOS/nixpkgs/nixos-20.09;

  outputs = { self, nixpkgs }: {

    defaultPackage.x86_64-linux =
      with import nixpkgs { system = "x86_64-linux"; };
      buildGo116Module rec {
        pname = "headscale";
        version = "0.0.1";
        src = self;
        vendorSha256 = "sha256-H3yvO9AvA80YfGFOapd2ByYvOiaXRsUTJcOe6MPkXkI=";
      };
  };
}
