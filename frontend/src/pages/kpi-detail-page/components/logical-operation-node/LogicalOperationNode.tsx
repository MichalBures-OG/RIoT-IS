import React, { useMemo } from 'react'
import { LogicalOperationNodeType } from '../editable-tree/EditableTree'
import EditableTreeNodeBase from '../editable-tree-node-base/EditableTreeNodeBase'
import { EffectFunction } from '../../../../util'

interface LogicalOperationNodeProps {
  type: LogicalOperationNodeType
  onClickHandler: EffectFunction
}

const LogicalOperationNode: React.FC<LogicalOperationNodeProps> = (props) => {
  const logicalOperationDenotation: string = useMemo(() => {
    switch (props.type) {
      case LogicalOperationNodeType.AND:
        return '∧ (AND)'
      case LogicalOperationNodeType.OR:
        return '∨ (OR)'
      case LogicalOperationNodeType.NOR:
        return '↓ (NOR)'
    }
  }, [props.type])

  return <EditableTreeNodeBase treeNodeContents={<p>{logicalOperationDenotation}</p>} onClickHandler={props.onClickHandler} />
}

export default LogicalOperationNode
