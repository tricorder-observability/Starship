import { useIntl } from '@umijs/max';
import { Form, Input } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      label={intl.formatMessage({
        id: 'code.eventSize',
      })}
      name="eventSize"
      rules={[
        { required: true, message: 'Please input eBPF event size!' },
        {
          validator: (rule, value) => {
            if (value === '0') {
              return Promise.resolve(true);
            }
            if (Number.isNaN(Number(value))) {
              return Promise.reject(false);
            }
            return Promise.resolve(true);
          },
          message: 'Please input number!',
        },
      ]}
    >
      <Input style={{ width: '100%' }} />
    </Form.Item>
  );
};
export default Index;
