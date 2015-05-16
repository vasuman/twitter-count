var TableControl = React.createClass({
    render: function() {
        return (
            <button className="control pure-button" onClick={ this.props.run } disabled={ !this.props.enabled }>
              <i className={ "fa fa-" + this.props.faIcon + (this.props.spin ? " fa-spin" : "") }></i>
            </button>
        );
    }
});

var DomainTable = React.createClass({
    getInitialState: function() {
        return { data: [], start: 0, requestPending: false }
    },
    loadFromServer: function(start) {
        var self = this;
        function success(data) {
            self.setState({
                start: start,
                data: data,
                requestPending: false
            });
        }
        function fail(data) {
            console.log(data);
            self.setState({ requestPending: false });
        }
        if (this.state.requestPending) {
            throw Exception("Request already pending");
        }
        this.setState({ requestPending: true });
        var params = {
            start: start,
            num: parseInt(this.props.size)
        };
        callAPI("listDomains", params, success, fail);
    },
    componentDidMount: function() {
        this._reset();
    },
    _reset: function() {
        this.setStart(0);
    },
    _refresh: function() {
        this.setStart(this.state.start);
    },
    _doNext: function() {
        var next = this.state.start + this.props.size;
        this.setStart(next);
    },
    _doPrev: function() {
        var prev = this.state.start - this.props.size;
        this.setStart(Math.max(prev, 0));
    },
    _canBack: function() {
        return !this.state.requestPending && this.state.start != 0;
    },
    _canNext: function() {
        return !this.state.requestPending && this.state.data.length == this.props.size;
    },
    setStart: function(start) {
        if (this._interval) {
            clearInterval(this._interval);
        }
        this.loadFromServer(start);
        var time = this.props.refreshInterval * 1000;
        this._interval = setInterval(this._refresh, time);
    },
    render: function() {
        var start = this.state.start + 1;
        function genDomainRow(item, idx) {
            return (
                <tr key={ item.name }>
                  <td> { start + idx } </td>
                  <td> <a href={ "domain?name=" + item.name }> { item.name } </a> </td>
                  <td> { item.count } </td>
                </tr>
            );
        }
        var nodes = this.state.data.map(genDomainRow),
            blankRow = <tr key="none" className="blank-row"> <td> { start } </td> <td colSpan="2"> - </td> </tr>;
        return (            
            <div className="table-container">
              <h1> Top Domains </h1>
              <table className="domains pure-table pure-table-striped">
                <thead>
                  <tr>
                    <th> # </th>
                    <th className="domain-name"> Name </th>
                    <th> Count </th>
                  </tr>
                </thead>
                <tbody>
                  { nodes.length != 0 ? nodes : blankRow }
                </tbody>
              </table>
              <div className="actions">
                <TableControl run={ this._reset }
                              name="Reset"
                              faIcon="step-backward"
                              enabled={ this._canBack() }/>
                <TableControl run={ this._doPrev }
                              name="Back"
                              faIcon="backward"
                              enabled={ this._canBack() }/>
                <TableControl run={ this._refresh }
                              name="Refresh"
                              faIcon="refresh"
                              enabled={ !this.state.requestPending }
                              spin={ this.state.requestPending }/>
                <TableControl run={ this._doNext }
                              name="Next"
                              faIcon="forward"
                              enabled={ this._canNext() }/>
              </div>
            </div>
        );
    }
});
