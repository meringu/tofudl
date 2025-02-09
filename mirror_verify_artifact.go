// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package tofudl

func (m *mirror) VerifyArtifact(artifactName string, artifactContents []byte, sumsFileContents []byte, signatureFileContent []byte) error {
	if m.pullThroughDownloader != nil {
		return m.pullThroughDownloader.VerifyArtifact(artifactName, artifactContents, sumsFileContents, signatureFileContent)
	}
	return verifyArtifact(m.keyRing, artifactName, artifactContents, sumsFileContents, signatureFileContent)
}
