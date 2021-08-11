local function cap_response(txn)
  res = txn.res
  -- this fails with
  --  runtime error: Cannot manipulate HAProxy channels in HTTP mode. from [C]: in method 'dup'
  -- print("response body = ", res:dup())
end

core.register_action("cap_response", { 'tcp-req', 'tcp-res', 'http-req', 'http-res' }, cap_response, 0)
