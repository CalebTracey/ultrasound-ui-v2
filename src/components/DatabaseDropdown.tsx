import React, { FC } from 'react'
import {
    Dropdown,
    DropdownToggle,
    DropdownMenu,
    DropdownItem,
    ButtonDropdown,
} from 'reactstrap'

interface Props {
    dropDownOpen: boolean
    dropDownToggle: () => void
    // databaseInitToggle: () => void
    databaseUpdateToggle: () => void
}

const DatabaseDropdown: FC<Props> = ({
    dropDownOpen,
    dropDownToggle,
    // databaseInitToggle,
    databaseUpdateToggle,
}) => (
    <Dropdown isOpen={dropDownOpen} toggle={dropDownToggle}>
        <DropdownToggle color="warning" caret>
            Admin
        </DropdownToggle>
        <DropdownMenu container="body">
            {/* <DropdownItem onClick={databaseInitToggle}>
<span className="span-text">Reset Database</span>
</DropdownItem> */}
            <DropdownItem onClick={databaseUpdateToggle}>
                <span className="span-text">Synchronize</span>
            </DropdownItem>
        </DropdownMenu>
    </Dropdown>
)

export default DatabaseDropdown
