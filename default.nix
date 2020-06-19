# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-wakeonlan-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-wakeonlan";

  modSha256 = "0nj0ny9692bqcw04fh74g8hqgfh3qc095fsq0y9cy677kp7l2q94";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
