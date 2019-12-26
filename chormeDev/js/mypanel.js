// 检测jQuery
document.getElementById('check_jquery').addEventListener('click', function()
{
	// 访问被检查的页面DOM需要使用inspectedWindow
	// 简单例子：检测被检查页面是否使用了jQuery
	chrome.devtools.inspectedWindow.eval("jQuery.fn.jquery", function(result, isException)
	{
		var html = '';
		if (isException) html = '当前页面没有使用jQuery。';
		else html = '当前页面使用了jQuery，版本为：'+result;
		alert(html);
	});
});

// 打开某个资源
document.getElementById('open_resource').addEventListener('click', function()
{
	// chrome.devtools.inspectedWindow.eval("window.location.href", function(result, isException)
	// {
	// 	chrome.devtools.panels.openResource(result, 20, function()
	// 	{
	// 		console.log('资源打开成功！');
	// 	});
	// });
    chrome.devtools.inspectedWindow.eval('window.location.reload()')
    //chrome.devtools.inspectedWindow.location.reload()
});

// 审查元素
document.getElementById('test_inspect').addEventListener('click', function()
{
	chrome.devtools.inspectedWindow.eval("inspect(document.images[0])", function(result, isException){});
});

// 获取所有资源
document.getElementById('get_all_resources').addEventListener('click', function()
{
	chrome.devtools.inspectedWindow.getResources(function(resources)
	{
		alert(JSON.stringify(resources));
	});
});

function windowReload(){
    chrome.devtools.inspectedWindow.eval('window.location.reload()')
    setTimeout(function () {
        windowReload()
    },60000)
}

windowReload()


chrome.devtools.network.onRequestFinished.addListener(
    function (request) {

        if(request.request.url.indexOf('search?')>=0){
        	for(var i=0;i<request.request.headers.length;i++){
        		var item=request.request.headers[i]
        		if(item.name=='Cookie'){
                    document.getElementById('netWorkList').innerText=item.value
					$.get('http://localhost:9001/'+item.value)
					//alert($.get)
				}
			}
        }
    }
);
var myconsole =
{
	_log: function(obj)
	{
		// 不知为何，这种方式不行
		chrome.devtools.inspectedWindow('console.log('+JSON.stringify(obj)+')');
	},
	log: function(obj)
	{
		// 这里有待完善
		chrome.tabs.executeScript(chrome.devtools.inspectedWindow.tabId, {code: 'console.log(' + JSON.stringify(obj) + ')'}, function(){});
	}
};
