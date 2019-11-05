import React, { PropTypes, Component } from 'react';
import '../App.css';
import { viewTypes, switchView } from '../actions';
import { connect } from 'react-redux';
import CommentView from '../components/CommentView';
import AccountView from '../components/AccountView';

class App extends Component {
  render() {
      /* let view;
       * if (this.props.view === viewTypes.COMMENT_VIEW) {
       *     view = (<CommentView onSwitchClick={this.props.onSwitchClick}/>);
       * } else {
       *     view = ();
       * }*/
      return (
      <div>
      <pre>
      {JSON.stringify(this.props)}
      </pre>
      <AccountView onSwitchClick={() => {this.props.onSwitchClick(); this.forceUpdate();}}/>
      </div>
          );

  }
}

App.defaultProps = {
    view: viewTypes.COMMENT_VIEW
};

App.propTypes = {
    view: PropTypes.string,
    onSwitchClick: PropTypes.func.Required
};

const mapStateToProps = (state, ownProps) => {
  return {
      view: state.view
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
      onSwitchClick: () => {
          let view;
          console.log(ownProps.view === viewTypes.COMMENT_VIEW);
          if (ownProps.view == viewTypes.COMMENT_VIEW) {
              view = viewTypes.ACCOUNT_VIEW;

          console.log("account");
          } else {
              view = viewTypes.COMMENT_VIEW
          console.log("comment");
          }
          console.log("Props" + JSON.stringify(ownProps));
          dispatch(switchView(view));
    }
  }
}

const AppC = connect(
    mapStateToProps,
    mapDispatchToProps
)(App)

export default AppC;
