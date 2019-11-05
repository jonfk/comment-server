import React, { PropTypes, Component } from 'react';

class AccountView extends Component {
    render() {
        return (
            <div>
                <p>
                    account view
                </p>
                <button onClick={this.props.onSwitchClick}>
                    switch
                </button>
            </div>
        );
    }
}

AccountView.propsTypes = {
    onSwitchClick: PropTypes.func.isRequired
}

export default AccountView;
