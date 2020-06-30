# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-wakeonlan-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-wakeonlan";

  modSha256 = "1gwm6zhxzc4nqmyhhn6p8cgvwj7dbcq1igafh0rsvhdfqngx3crd";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
