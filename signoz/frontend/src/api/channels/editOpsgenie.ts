import axios from 'api';
import { ErrorResponseHandler } from 'api/ErrorResponseHandler';
import { AxiosError } from 'axios';
import { ErrorResponse, SuccessResponse } from 'types/api';
import { PayloadProps, Props } from 'types/api/channels/editOpsgenie';

const editOpsgenie = async (
	props: Props,
): Promise<SuccessResponse<PayloadProps> | ErrorResponse> => {
	try {
		const response = await axios.put(`/channels/${props.id}`, {
			name: props.name,
			opsgenie_configs: [
				{
					send_resolved: true,
					api_key: props.api_key,
					description: props.description,
					priority: props.priority,
					message: props.message,
					details: {
						...props.detailsArray,
					},
				},
			],
		});

		return {
			statusCode: 200,
			error: null,
			message: 'Success',
			payload: response.data.data,
		};
	} catch (error) {
		return ErrorResponseHandler(error as AxiosError);
	}
};

export default editOpsgenie;
