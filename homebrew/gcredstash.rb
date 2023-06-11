require "formula"

class Gcredstash < Formula
  VERSION = "0.4.2"

  desc "Manages credentials using AWS Key Management Service (KMS) and DynamoDB"
  homepage "https://github.com/kgaughan/gcredstash"
  version VERSION

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kgaughan/gcredstash/releases/download/v#{VERSION}/gcredstash_#{VERSION}_darwin_x86_64.tar.xz"
      sha256 "8f1cae06fd80c9ce5d553150b7309d42ed3bb0de293c09a87779b45603ed2e50"

      def install
        bin.install "gcredstash"
      end
    elsif Hardware::CPU.arm?
      url "https://github.com/kgaughan/gcredstash/releases/download/v#{VERSION}/gcredstash_#{VERSION}_darwin_arm64.tar.xz"
      sha256 "e5d51249da470af6679c9b79d71850caa6d4d300ac9d0939c3d5d720b43f8d29"

      def install
        bin.install "gcredstash"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/kgaughan/gcredstash/releases/download/v#{VERSION}/gcredstash_#{VERSION}_linux_x86_64.tar.xz"
      sha256 "29161eab1437e83c550fbb10ffef7ea1ceaa95872668a3d76e60262e9203cb06"

      def install
        bin.install "gcredstash"
      end
    elsif Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kgaughan/gcredstash/releases/download/v#{VERSION}/gcredstash_#{VERSION}_linux_arm64.tar.xz"
      sha256 "d9082265289ed68fb70367ca24e9902e6037a02635c4ea3dd6c3cf4de3902892"

      def install
        bin.install "gcredstash"
      end
    end
  end
end
