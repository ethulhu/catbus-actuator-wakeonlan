# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-wakeonlan-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-wakeonlan";

  modSha256 = "1vv9g9g55zpq68snpk3m6ashzwipy5w3xng06l7i5pbjw1n27m0g";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
