import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { useIntl } from '@umijs/max';
import { Button, Form, Input, Select, Space } from 'antd';

type IProps = {};

const Index = (props: IProps) => {
  const intl = useIntl();
  return (
    <Form.Item
      wrapperCol={{ offset: 0, span: 15 }}
      label={intl.formatMessage({
        id: 'code.collector',
      })}
      required={true}
    >
      <Form.List
        name="schemaAttr"
        initialValue={[
          {
            name: null,
            type: null,
          },
        ]}
      >
        {(fields, { add, remove }) => (
          <>
            {fields.map(({ key, name, ...restField }) => (
              <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                <Form.Item
                  {...restField}
                  name={[name, 'name']}
                  rules={[{ required: true, message: 'attibute' }]}
                >
                  <Input placeholder="attibute name" />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'type']}
                  rules={[{ required: true, message: 'attibute type' }]}
                >
                  <Select
                    placeholder="attibute type"
                    style={{
                      width: 166,
                    }}
                  >
                    <Select.Option value={0}>bool</Select.Option>
                    <Select.Option value={1}>date</Select.Option>
                    <Select.Option value={2}>int</Select.Option>
                    <Select.Option value={3}>integer</Select.Option>
                    <Select.Option value={4}>json</Select.Option>
                    <Select.Option value={5}>jsonb</Select.Option>
                    <Select.Option value={6}>text</Select.Option>
                  </Select>
                </Form.Item>
                {fields.length > 1 && <MinusCircleOutlined onClick={() => remove(name)} />}
              </Space>
            ))}
            <Form.Item>
              <Button
                type="dashed"
                onClick={() => add()}
                block
                icon={<PlusOutlined />}
                style={{ width: '100%' }}
              >
                {intl.formatMessage({
                  id: 'code.addField',
                })}
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
    </Form.Item>
  );
};
export default Index;
