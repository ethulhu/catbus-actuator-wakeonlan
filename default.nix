# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-wakeonlan-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-wakeonlan";

  modSha256 = "0j3yd7da2gicbkwxf27sshdd32qzkjvj33pfhljafizl7fxyd7lq";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
