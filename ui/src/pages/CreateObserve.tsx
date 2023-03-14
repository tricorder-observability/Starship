import { MinusCircleOutlined, PlusOutlined, UploadOutlined } from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-components';
import { history, useIntl } from '@umijs/max';
import { Button, Card, Form, Input, message, Select, Space, Upload } from 'antd';
import React, { useEffect, useState } from 'react';
import { createModule } from '../services/ant-design-pro/api';

const width = '100%';
const Code: React.FC = () => {
  const [form] = Form.useForm();
  const intl = useIntl();
  const [fileContent, setFileContent] = useState<any>([]);
  const onFinish = async (values: any) => {
    try {
      const params = {
        ebpf: {
          code: values.code,
          fmt: 0,
          lang: 0,
          probes: values.probes,
        },
        name: values.name,
        wasm: {
          code: fileContent,
          // TODO(zhoujie): fmt default value
          fmt: 0,
          fn_name: values.fn,
          // TODO(zhoujie): lang default value
          lang: 0,
          output_schema: {
            fields: values.schemaAttr,
          },
        },
      };
      const msg = await createModule(params);
      if (msg.code === 200) {
        message.success('success');
        history.push('/module-list');
        sessionStorage.setItem('codeCache', '');
      } else {
        message.error(msg.message);
      }
      return msg.data;
    } catch (error) {}
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('failed:', errorInfo);
  };

  useEffect(() => {
    const unlisten = history.listen(() => {
      // Every time the route changes, it will go here
      sessionStorage.setItem('codeCache', JSON.stringify(form.getFieldsValue()));
    });
    const beforeunloadCallback = (event: any = window.event) => {
      // This api has been abolished. It is only done for compatibility with some browsers
      event.returnValue = 'You might lost changes on this page';
      return event.returnValue;
    };
    // The beforeupload event will only be triggered if the user has mouse interaction on the current page
    window.addEventListener('beforeunload', beforeunloadCallback);
    return () => {
      unlisten();
      window.removeEventListener('beforeunload', beforeunloadCallback);
    };
  }, [form]);

  const readFileContent = (info: any) => {
    if (info.file.status === 'done') {
      const reader = new FileReader();
      reader.onload = (e) => {
        const arrBuffer: any = e.target?.result;
        const uint8Array = new Uint8Array(arrBuffer);
        setFileContent(Array.from(uint8Array));
      };
      reader.readAsArrayBuffer(info.file.originFileObj);
    }
  };

  return (
    <PageContainer>
      <Card
        style={{
          borderRadius: 8,
        }}
        bodyStyle={{
          backgroundImage:
            'radial-gradient(circle at 97% 10%, #EBF2FF 0%, #F5F8FF 28%, #EBF1FF 124%)',
        }}
      >
        <Form
          name="basic"
          labelCol={{ span: 5 }}
          wrapperCol={{ span: 15 }}
          initialValues={
            sessionStorage.getItem('codeCache')
              ? JSON.parse(sessionStorage.getItem('codeCache') || '{}')
              : {}
          }
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          autoComplete="off"
          form={form}
        >
          <Form.Item
            label={intl.formatMessage({
              id: 'code.name',
            })}
            name="name"
            rules={[
              { required: true, message: 'Please input name!' },
              { type: 'string', max: 50, message: 'Up to 50 characters!' },
            ]}
          >
            <Input style={{ width: width }} />
          </Form.Item>
          <Form.Item
            label={intl.formatMessage({
              id: 'code.code',
            })}
            name="code"
            rules={[{ required: true, message: 'Please input ebpf!' }]}
          >
            <Input.TextArea rows={8} style={{ width: width }} />
          </Form.Item>
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
            <Input style={{ width: width }} />
          </Form.Item>
          <Form.Item
            label={intl.formatMessage({
              id: 'code.perfBuffers',
            })}
            name="perfBuffers"
            rules={[{ required: true, message: 'Please input eBPF perf buffers!' }]}
          >
            <Input style={{ width: width }} />
          </Form.Item>
          <Form.Item
            wrapperCol={{ offset: 0, span: 15 }}
            label={intl.formatMessage({
              id: 'code.kprobe',
            })}
            required={true}
          >
            <Form.List
              name="probes"
              initialValue={[
                {
                  target: null,
                  entry: null,
                  return: null,
                },
              ]}
            >
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, ...restField }) => (
                    <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                      <Form.Item
                        {...restField}
                        name={[name, 'target']}
                        rules={[{ required: true, message: 'target' }]}
                      >
                        <Input placeholder="target" />
                      </Form.Item>
                      <Form.Item
                        {...restField}
                        name={[name, 'entry']}
                        rules={[{ required: true, message: 'entry' }]}
                      >
                        <Input placeholder="entry" />
                      </Form.Item>
                      <Form.Item
                        {...restField}
                        name={[name, 'return']}
                        rules={[{ required: true, message: 'return' }]}
                      >
                        <Input placeholder="return" />
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
                      style={{ width: width }}
                    >
                      {intl.formatMessage({
                        id: 'code.addKprobe',
                      })}
                    </Button>
                  </Form.Item>
                </>
              )}
            </Form.List>
          </Form.Item>
          <Form.Item
            name="wasm"
            label={intl.formatMessage({
              id: 'code.wasm',
            })}
            extra="WASM byte code(.wasm .wat)"
            rules={[{ required: true, message: 'Please input wasm code!' }]}
          >
            <Upload
              action=""
              maxCount={1}
              accept=".wasm,.wat"
              showUploadList={false}
              onChange={readFileContent}
            >
              <Button icon={<UploadOutlined />}>Click to upload</Button>
            </Upload>
          </Form.Item>
          <Form.Item
            label={intl.formatMessage({
              id: 'code.fn',
            })}
            name="fn"
            rules={[{ required: true, message: 'Please input function name!' }]}
          >
            <Input style={{ width: width }} />
          </Form.Item>
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
                      style={{ width: width }}
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
          <Form.Item wrapperCol={{ offset: 5, span: 16 }}>
            <Button type="primary" htmlType="submit">
              {intl.formatMessage({
                id: 'button.submit',
              })}
            </Button>
            <Button
              type="primary"
              onClick={() => {
                sessionStorage.setItem('codeCache', '');
                form.setFieldsValue({
                  name: null,
                  code: null,
                  eventSize: null,
                  perfBuffers: null,
                  probes: null,
                  fn: null,
                  wasm: null,
                  schemaName: null,
                  schemaAttr: null,
                });
              }}
              style={{
                marginLeft: 10,
              }}
            >
              {intl.formatMessage({
                id: 'button.clear',
              })}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </PageContainer>
  );
};

export default Code;
