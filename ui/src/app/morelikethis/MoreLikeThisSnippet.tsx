import * as React from "react";
import {Card, CardBody, CardFooter} from "reactstrap";
import './MoreLikeThisSnippet.css';
import {MoreLikeThisOption} from "./MoreLikeThisOption";
import ReadableId from "../common/ReadableId";

interface MoreLikeThisSnippetProps {
    id: string,
    text: string
    onSelect: ((id: string, relevant: MoreLikeThisOption) => void)
    relevant: MoreLikeThisOption
}

class MoreLikeThisSnippet extends React.Component<MoreLikeThisSnippetProps, any> {
    constructor(props: any, context: any) {
        super(props, context);
        this.state = {};
    }

    onClickYes = () => {
        this.props.onSelect(this.props.id, MoreLikeThisOption.YES);
    };

    onClickNo = () => {
        this.props.onSelect(this.props.id, MoreLikeThisOption.NO);
    };

    onClickMaybe = () => {
        this.props.onSelect(this.props.id, MoreLikeThisOption.MAYBE);
    };

    render() {
        return (
            <Card className="more-like-this-snippet">
                <CardBody>
                    <span className="small"><strong><ReadableId id={this.props.id}/></strong></span>
                    <p className="small snippet-text">{this.props.text}</p>
                </CardBody>
                <CardFooter>
                    <span>Is dit fragment relevant?</span>
                    <div className="float-right btn-group btn-group-sm" role="group">
                        <button type="button"
                                className={`btn btn-outline-success ${this.props.relevant === MoreLikeThisOption.YES ? "active" : ""}`}
                                onClick={this.onClickYes}>
                            <i className='fas fa-check-circle'/>
                            &nbsp;
                            Ja
                        </button>
                        <button type="button"
                                className={`btn btn-outline-danger ${this.props.relevant === MoreLikeThisOption.NO ? "active" : ""}`}
                                onClick={this.onClickNo}>
                            <i className='fas fa-times-circle'/>
                            &nbsp;
                            Nee
                        </button>
                        <button type="button"
                                className={`btn btn-outline-secondary ${this.props.relevant === MoreLikeThisOption.MAYBE ? "active" : ""}`}
                                onClick={this.onClickMaybe}>
                            <i className='fas fa-question-circle'/>
                            &nbsp;
                            Blanco
                        </button>
                    </div>
                </CardFooter>
            </Card>
        );
    }
}

export default MoreLikeThisSnippet;