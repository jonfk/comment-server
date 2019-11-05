import React, { PropTypes, Component } from 'react';

class CommentView extends Component {
    render() {
        return (
            <div>
                <p>
                    comment view
                </p>
                <button onClick={this.props.onSwitchClick}>
                    switch
                </button>
            </div>
        );
    }
}

CommentView.propsTypes = {
    onSwitchClick: PropTypes.func.isRequired
}

export default CommentView;
