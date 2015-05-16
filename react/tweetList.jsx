var TweetCard = React.createClass({
    render: function() {
        var tweet = this.props.tweet,
            userLink = "https://twitter.com/" + tweet.userHandle,
            tweetLink = userLink + "/status/" + tweet.id;
        return (
            <div className="tweet-card">
              <span className="user-handle">
                <a className="userLink" href={ userLink }>
                  @{ tweet.userHandle }
                </a>
              </span>
              <a href={ tweetLink }>...</a>
            </div>
        );
    }
});

var TweetList = React.createClass({
    getInitialState: function() {
        return { tweets: [], loading: false }
    },

    componentDidMount: function() {
        this.nextBatch();
    },
    
    nextBatch: function() {
        var self = this,
            params = {
                domain: this.props.domainName,
                idx: this.state.tweets.length,
                num: this.props.batchSize
            }
        function success(data) {
            console.log(data);
            var newTweets = self.state.tweets.concat(data)
                self.setState({
                    loading: false,
                    tweets: newTweets,
                });
        }
        function fail(data) {
            console.log(data);
            self.setState({ loading: false });
        }
        this.setState({ loading: true });
        callAPI("getTweets", params, success, fail);
    },
    
    render: function() {
        var nodes = this.state.tweets.map(function(tweet) {
            console.log(tweet);
            return (
                <TweetCard key={ tweet.id  } tweet={ tweet }/>
            );
        });
        var end;
        if (!this.state.isLoading) {
            end = (
                <div className="list-end tweet-card" onClick={ this.nextBatch }>
                  More
                </div>
            );
        } else {
            end = (
                <div className="list-end tweet-card">
                  <i className="fa fa-2x fa-spinner fa-spin"></i>
                </div>
            );
        }
        return (
            <div className="tweet-list">
              <div className="list-start tweet-card">
                <i className="fa fa-2x fa-twitter twitter-icon"></i>
                for
                <span className="domain-name">
                  { this.props.domainName }
                </span>
              </div>
              { nodes }
              { end }
            </div>
        );
    }
});
