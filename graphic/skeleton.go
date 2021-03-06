// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	)

// Skeleton contains armature information.
type Skeleton struct {
	core.INode
	inverseBindMatrices []math32.Matrix4
	boneMatrices        []math32.Matrix4
	bones               []*core.Node
}

// NewSkeleton creates and returns a pointer to a new Skeleton.
func NewSkeleton(node core.INode) *Skeleton {

	sk := new(Skeleton)
	sk.INode = node
	sk.boneMatrices = make([]math32.Matrix4, 0)
	sk.bones = make([]*core.Node, 0)
	return sk
}

// AddBone adds a bone to the skeleton along with an optional inverseBindMatrix.
func (sk *Skeleton) AddBone(node *core.Node, inverseBindMatrix *math32.Matrix4) {

	// Useful for debugging:
	//node.Add(NewAxisHelper(0.2))

	sk.bones = append(sk.bones, node)
	sk.boneMatrices = append(sk.boneMatrices, *math32.NewMatrix4())
	if inverseBindMatrix == nil {
		inverseBindMatrix = math32.NewMatrix4() // Identity matrix
	}

	sk.inverseBindMatrices = append(sk.inverseBindMatrices, *inverseBindMatrix)
}

// Bones returns the list of bones in the skeleton.
func (sk *Skeleton) Bones() []*core.Node {

	return sk.bones
}

// BoneMatrices calculates and returns the bone world matrices to be sent to the shader.
func (sk *Skeleton) BoneMatrices() []math32.Matrix4 {

	// Obtain inverse matrix world
	var invMat math32.Matrix4
	node := sk.GetNode()
	node.UpdateMatrixWorld()
	nMW := node.MatrixWorld()
	err := invMat.GetInverse(&nMW)
	if err != nil {
		log.Error("Skeleton.BoneMatrices: inverting matrix failed!")
	}

	// Update bone matrices
	for i := range sk.bones {
		sk.bones[i].UpdateMatrixWorld()
		bMat := sk.bones[i].MatrixWorld()
		bMat.MultiplyMatrices(&bMat, &sk.inverseBindMatrices[i])
		sk.boneMatrices[i].MultiplyMatrices(&invMat, &bMat)
	}

	return sk.boneMatrices
}
